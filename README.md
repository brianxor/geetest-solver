# Geetest Solver

A blazing-fast and ready-to-use API to solve Geetest captchas.

It uses no external dependencies, such as OpenCV, saving you from the hassle of installation. The solution relies solely on standard libraries, making it lightweight and easy to use.

## Installation

```sh
git clone https://github.com/tylerdevx/geetest-solver.git
cd geetest-solver
go run .
```

## Usage

### Base URL

```
http://127.0.0.1:8080
```

### V4 Puzzle

Recognition Time: **<10ms**

Total time depends on the speed of your proxy for the request.

#### Endpoint

```
POST /geetest/v4/puzzle/solve
```

#### Request Body

```json
{
  "websiteUrl": "https://example.com/",
  "captchaId": "...",
  "userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36",
  "proxy": "http://user:pass@ip:port"
}
```

- The `proxy` field is **optional**.

## Upcoming Features

- Support for additional Geetest modes
- Performance improvements and optimizations
- More API customization options

Stay tuned!

---

If you appreciate my work, please leave a star ⭐️ on GitHub! 🙏
