package domain

type ServerModel interface {
	GetScheme() string
	GetPort() string
	GetAddress() string
	GetCertFilePath() string
	GetKeyFilePath() string
}

type Server struct {
	Scheme   string
	Host     string
	Port     string
	CertPath string
	KeyPath  string
}

func (s *Server) GetScheme() string {
	return s.Scheme
}

func (s *Server) GetPort() string {
	return s.Port
}

func (s *Server) GetAddress() string {
	return s.Host + ":" + s.Port
}

func (s *Server) GetCertFilePath() string {
	return s.CertPath
}

func (s *Server) GetKeyFilePath() string {
	return s.KeyPath
}

type Servers []Server

func (s *Servers) Next() ServerModel {
	if len(*s) == 0 {
		return nil
	}
	server := (*s)[0]
	*s = (*s)[1:]
	return &server
}
