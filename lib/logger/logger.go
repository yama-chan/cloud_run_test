package logger

type Logger interface {
	Debug(message string) //1
	Info(message string)  //2
	Warn(message string)  //3
	Error(message string) //4
	Fatal(message string) //5
	Panic(message string) //6
}
