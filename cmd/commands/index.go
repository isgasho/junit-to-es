package commands

import (
	"os"

	"github.com/rs/zerolog"

	"github.com/spf13/cobra"
)

var (
	indexSetup     bool
	crawlerWorkers int
)

func init() {
	rootCmd.AddCommand(indexCmd)

}

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index test results from maven surefire plugin into Elasticsearch",
	Run: func(cmd *cobra.Command, args []string) {
		crawler := Crawler{
			log: zerolog.New(zerolog.NewConsoleWriter{Out: os.Stdout}).
			Level(func() zerolog.Level{
				if os.Getenv("DEBUG") != "" {
					return zerolog.DebugLevel
				} else {
					return zerolog.InfoLevel
				}
			}()).With().Timestamp().Logger(),
			workers: crawlerWorkers,
			queue: make(chan string, crawlerWorkers),

		},
	},
}


func (c *Crawler) setupIndex() error {
	mapping := `{
    "mappings": {
      "_doc": {
        "properties": {
          "id":         { "type": "keyword" },
          "image_url":  { "type": "keyword" },
          "title":      { "type": "text", "analyzer": "english" },
          "alt":        { "type": "text", "analyzer": "english" },
          "transcript": { "type": "text", "analyzer": "english" },
          "published":  { "type": "date" },
          "link":       { "type": "keyword" },
          "news":       { "type": "text", "analyzer": "english" }
        }
      }
    }
		}`
	return c.store.CreateIndex(mapping)
}
