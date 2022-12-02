package clients

type HttpClient struct {
	EmployeeClient   *EmployeeClient
	DigitalClient    *DigitalClient
	MasterDataClient *MasterDataClient
	ContractClient   *ContractClient
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		EmployeeClient:   NewEmployeeClient(),
		DigitalClient:    NewDigitalClient(),
		MasterDataClient: NewMasterDataClient(),
		ContractClient:   NewContractClient(),
	}
}
