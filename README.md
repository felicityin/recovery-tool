# Description

A tool for recovering threshold private keys.

# Install dependencies

You'll need to have the following tools installed on your machine.

- [go](https://golang.org/)

# Build

```
git clone https://github.com/felicityin/recovery-tool.git
cd recovery-tool
go build
```

# Usage

## Recover keys

```
./recovery-tool recover -i input.yaml -o output.yaml
```

## Get balance

```
# Solana
./recovery-tool balance -addr HXx8Ky1aY7GBLUghbadKais5QHdeJfdQ7mmgR9j4sqNK

# Solana USDT
./recovery-tool balance -addr HXx8Ky1aY7GBLUghbadKais5QHdeJfdQ7mmgR9j4sqNK -coin Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB

# Aptos
./recovery-tool balance -chain apt -addr 68c709c6614e29f401b6bfdd0b89578381ef0fb719515c03b73cf13e45550e06

# Polkadot
./recovery-tool balance -chain dot -addr 1hXFGPkKRcXQH4FrrRA7tVvfw1Ghaoy5wsW5TMjTokq1QPk
```

## Transfer

```
# Solana
./recovery-tool transfer -fromkey 078fe2333b309a95f8bc59f6e03a10c4b7b51f3e12b7ccd4a62c41363a08437a -to FUQ3cTZpuB23cohYUFWTfnK6AHTEKZ9u5vAbkBGdTFdD -amount 0.00001

# Solana USDT
./recovery-tool transfer -fromkey 078fe2333b309a95f8bc59f6e03a10c4b7b51f3e12b7ccd4a62c41363a08437a -to FUQ3cTZpuB23cohYUFWTfnK6AHTEKZ9u5vAbkBGdTFdD -amount 0.001 -coin Es9vMFrzaCERmJfrF4H2FYD4KCoNkY11McCe8BenwNYB

# Aptos
./recovery-tool transfer -chain apt -fromkey 0f2020d5f3ff4a08919d6e5f9058c47946dffaa620ea10e0d884c078dfa6ba23 -to 891b2e8cf70a98759435a2efed6636d7b054964ab2c55e7d51ef1ab3e32850a0 -amount 0.00001

# Polkadot
./recovery-tool transfer -chain dot -fromkey 0d9e11aeaa5d1f00565386799fa6e04e51c3b8087113e972d0dfc4bcc26ad9dc -to 15Yo7C1g4YLn3MZwdZo6Gwp9SMcXqUpUQXVkZRGsgaKJoUV5 -amount 0.00001
```
