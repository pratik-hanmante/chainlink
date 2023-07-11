package mercury

import (
	"encoding/json"

	"github.com/pkg/errors"

	libocr2 "github.com/smartcontractkit/libocr/offchainreporting2plus"

	relaymercuryv0 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v0"
	relaymercuryv1 "github.com/smartcontractkit/chainlink-relay/pkg/reportingplugins/mercury/v1"
	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"

	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/mercury/config"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury"
	mercuryv0 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v0"
	mercuryv1 "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/mercury/v1"
)

type Config interface {
	MaxSuccessfulRuns() uint64
}

func NewServices(
	jb job.Job,
	ocr2Provider relaytypes.MercuryProvider,
	pipelineRunner pipeline.Runner,
	runResults chan pipeline.Run,
	lggr logger.Logger,
	argsNoPlugin libocr2.MercuryOracleArgs,
	cfg Config,
	chEnhancedTelem chan ocrcommon.EnhancedTelemetryMercuryData,
	chainHeadTracker mercury.ChainHeadTracker,
) ([]job.ServiceCtx, error) {
	if jb.PipelineSpec == nil {
		return nil, errors.New("expected job to have a non-nil PipelineSpec")
	}

	var pluginConfig config.PluginConfig
	err := json.Unmarshal(jb.OCR2OracleSpec.PluginConfig.Bytes(), &pluginConfig)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	err = config.ValidatePluginConfig(pluginConfig)
	if err != nil {
		return nil, err
	}
	lggr = lggr.Named("MercuryPlugin").With("jobID", jb.ID, "jobName", jb.Name.ValueOrZero())

	switch ocr2Provider.ReportSchemaVersion() {
	case 0:
		ds := mercuryv0.NewDataSource(
			pipelineRunner,
			jb,
			*jb.PipelineSpec,
			lggr,
			runResults,
			chEnhancedTelem,
			chainHeadTracker,
			ocr2Provider.ContractTransmitter(),
			pluginConfig.InitialBlockNumber.Ptr(),
		)
		argsNoPlugin.MercuryPluginFactory = relaymercuryv0.NewFactory(
			ds,
			lggr,
			ocr2Provider.OnchainConfigCodec(),
			ocr2Provider.ReportCodecV0(),
		)
	case 1:
		ds := mercuryv1.NewDataSource(
			pipelineRunner,
			jb,
			*jb.PipelineSpec,
			lggr,
			runResults,
			chEnhancedTelem,
		)
		argsNoPlugin.MercuryPluginFactory = relaymercuryv1.NewFactory(
			ds,
			lggr,
			ocr2Provider.OnchainConfigCodec(),
			ocr2Provider.ReportCodecV1(),
		)
	default:
		panic("unknown schema version")
	}

	oracle, err := libocr2.NewOracle(argsNoPlugin)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return []job.ServiceCtx{ocr2Provider, ocrcommon.NewResultRunSaver(
		runResults,
		pipelineRunner,
		make(chan struct{}),
		lggr,
		cfg.MaxSuccessfulRuns(),
	),
		job.NewServiceAdapter(oracle)}, nil
}
