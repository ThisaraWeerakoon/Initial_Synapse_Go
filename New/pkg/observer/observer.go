package observer



//Listen and update the log level
type Observer interface{
	UpdateLogger()
}