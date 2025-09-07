package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/rafaellevissa/rox-partner/internal/app"
)

func NewRootCmd() *cobra.Command {
	var cfgFile string

	rootCmd := &cobra.Command{
		Use:   "b3ingestor",
		Short: "B3 Ingestor - Process and query B3 trading data",
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./configs/config.yaml)")

	cobra.OnInitialize(func() {
		if cfgFile != "" {
			viper.SetConfigFile(cfgFile)
		} else {
			viper.SetConfigName("config")
			viper.AddConfigPath("./configs")
		}
		viper.AutomaticEnv()
		_ = viper.ReadInConfig()
	})

	rootCmd.AddCommand(newIngestCmd())

	return rootCmd
}

func newIngestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ingest [zipfile|dir]",
		Short: "Ingest trades from a B3 zip file or directory of zip files",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dsn := viper.GetString("database.dsn")
			tmp := viper.GetString("ingestion.tmp_dir")
			batch := viper.GetInt("ingestion.batch_size")
			return app.IngestTrades(dsn, args[0], tmp, batch)
		},
	}
}
