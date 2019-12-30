package source

// InjectNewHTTPClient replaces mock HTTPClient for testing. Use the function in only test case.
func InjectNewHTTPClient(c httpClient) { newHTTPClient = func() httpClient { return c } }

// FixNewHTTPClient fixes HTTPClient constructor with original one. Use the function in only test case.
func FixNewHTTPClient() { newHTTPClient = newNormalHTTPClient }
