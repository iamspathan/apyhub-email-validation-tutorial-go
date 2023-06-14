## Getting Started

Once you created the account, Copy your [apy-token] and paste in [main.go](https://github.com/iamspathan/apyhub-email-validation-tutorial-go/blob/master/main.go):

```
func validateEmail(email string) (bool, error) {
	url := "https://api.apyhub.com/validate/email/dns"

	payload := struct {
		Email string `json:"email"`
	}{
		Email: email,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return false, err
	}

	if err != nil {
		return false, err
	}
	// Add your apy-token here.
	
	req.Header.Add("apy-token", "*********** ADD YOUR SECRET APY TOKEN HERE **************")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	var response struct {
		Valid bool `json:"data"`
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return false, err
	}

	return response.Valid, nil
}
```

Run the development server:

```
go run main.go
```

You will see the below message in console.

```
Server started on port 3000

```

Navigate to [http://localhost:3000](http://localhost:3000)

