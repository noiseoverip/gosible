package transport

type Transport interface {
	Exec(command string, args... string) (resultCode int, stdout string, stderr string)
}
