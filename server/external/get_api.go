package external

type getAPI interface {
	getCryptocurrency() ([]byte, error)
	getFiat() ([]byte, error)
}
