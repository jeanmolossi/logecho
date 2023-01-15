package logecho

type StdLogger interface {
	WithFields(f Fields) StdLogger
	Println(i ...interface{})
	Printf(format string, args ...interface{})
	Debugln(i ...interface{})
	Debugf(format string, args ...interface{})
	Infoln(i ...interface{})
	Infof(format string, args ...interface{})
	Warnln(i ...interface{})
	Warnf(format string, args ...interface{})
	Errorln(i ...interface{})
	Errorf(format string, args ...interface{})
	Fatalln(i ...interface{})
	Fatalf(format string, args ...interface{})
	Panicln(i ...interface{})
	Panicf(format string, args ...interface{})
}
