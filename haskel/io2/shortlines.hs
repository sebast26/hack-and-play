main = do
    contents <- getContents
    putStr (shortLinesOnly contents)

shortLinesOnly input =
    let allLines = lines input
        shortLines = filter (\line -> length line < 10) allLines
        result = unlines shortLines
    in result