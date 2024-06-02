# GO-TAMBOON

Uses go version `1.22.3`

Steps to get started:
1. Clone repo `git clone <repo urls>`
2. Get omise public & private keys from [here](https://dashboard.omise.co/v2/settings/keys). **Esure** you are in test mode
3. Run `cp env.example .env` and add your keys to the env
4. Or you can `export <KEY_NAME>=<value>`
5. Then simply run `go run cmd/tamboon.go ./data/fng.1000.csv.rot128`

Project Structure:



