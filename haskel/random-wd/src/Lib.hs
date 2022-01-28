module Lib
    ( someFunc
    ) where

import System.Random
import Control.Monad(when)

someFunc :: IO ()
someFunc = do
    print $ threeCoins (mkStdGen 21) 
    gen <- getStdGen
    print $ take 20 (randomRs ('a', 'z') gen)
    print $ take 5 $ (randoms (mkStdGen 11) :: [Int])
    askForNumber gen 

threeCoins :: StdGen -> (Bool, Bool, Bool)
threeCoins gen =
    let (firstCoin, newGen) = random gen
        (secondCoin, newGen') = random newGen
        (thirdCoin, newGen'') = random newGen'
    in (firstCoin, secondCoin, thirdCoin)

askForNumber :: StdGen -> IO ()  
askForNumber gen = do  
    let (randNumber, newGen) = randomR (1,10) gen :: (Int, StdGen)  
    putStr "Which number in the range from 1 to 10 am I thinking of? "  
    numberString <- getLine  
    when (not $ null numberString) $ do  
        let number = read numberString  
        if randNumber == number   
            then putStrLn "You are correct!"  
            else putStrLn $ "Sorry, it was " ++ show randNumber  
        askForNumber newGen  
