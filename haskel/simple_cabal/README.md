# simple_cabal

To run simple `hello.hs` file:

```bash
runghc hello.hs
```

## Cabal

To initialize project `cabal init -i`.

To run project `cabal run`.

To add dependency `scotty >= 0.12`, you need to edit cabal file and adjust lines:
```bash
build-depends:    
    base ^>=4.17.0.0,
    scotty >= 0.12
```

To use `Scotty` update `app/Main.hs` with:
```haskell
{-# LANGUAGE OverloadedStrings #-}
import Web.Scotty

import Data.Monoid (mconcat)

main = scotty 3000 $
    get "/:word" $ do
        word <- param "word"
        html $ mconcat ["<h1>", word, " World!</h1>"]
```

## Repl

To use repl use `cabal repl`.

## To lock dependencies

`cabal freeze`

