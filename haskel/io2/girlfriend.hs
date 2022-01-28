import System.IO

main = do
    withFile' "girlfriend.txt" ReadMode (\handle -> do
        contents <- hGetContents handle
        putStr contents)

-- for build in function use withFile
withFile' :: FilePath -> IOMode -> (Handle -> IO a) -> IO a
withFile' path mode f = do
    handle <- openFile path mode
    result <- f handle
    hClose handle
    return result


-- main = do
    -- contents <- readFile "girlfriend.txt"
    -- putStr contents

-- main = do
    -- contents <- readFile "girlfriend.txt"
    -- writeFile "girlfriendcaps.txt" (map toUpper contents)