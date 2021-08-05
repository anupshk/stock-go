package cli

func GetAllClients() []Client {
	var clients []Client
	err := DB.Model(&clients).Select()
	if err != nil {
		panic(err)
	}
	return clients
}

func InsertClient(client *Client) (interface{}, error) {
	return DB.Model(client).Insert()
}
