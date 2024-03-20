package justdb

type Options struct {
	// DirPath is the directory path where the data will be stored
	DirPath string
	// AutoMergeCronExpr is the cron expression for auto merge
	AutoMergeCronExpr string
	// ActiveDataFileMaxSize is the maximum size of the active data file
	ActiveDataFileMaxSize int64
}
