import System.Environment
import System.Directory
import System.IO
import Data.List
import Control.Exception

main = do
  (command:argList) <- getArgs
  dispatch command argList

dispatch :: String -> [String] -> IO ()
dispatch "view" = view
dispatch "add" = add
dispatch "remove" = remove
dispatch command = doesntExist command

add :: [String] -> IO ()
add [fileName, todoItem] = appendFile fileName (todoItem ++ "\n")
add _ = putStrLn "The add command takes exactly two arguments"

view :: [String] -> IO ()
view [fileName] = do
  contents <- readFile fileName
  let todoTasks = lines contents
      numberedTasks = zipWith (\n line -> show n ++ " - " ++ line) [0..] todoTasks
  putStr $ unlines numberedTasks
view _ = putStrLn "the view command takes exactly one argument"

remove :: [String] -> IO ()
remove [fileName, numberString] = do
  contents <- readFile fileName
  let todoTasks = lines contents
      number = read numberString
      newTodoItems = unlines $ delete (todoTasks !! number) todoTasks
      numberedTasks = zipWith (\n line -> show n ++ " - " ++ line) [0..] $ lines newTodoItems
  bracketOnError (openTempFile "." "temp")
    (\(tempName, tempHandle) -> do
      hClose tempHandle
      removeFile tempName)
    (\(tempName, tempHandle) -> do
      hPutStr tempHandle newTodoItems
      hClose tempHandle
      removeFile fileName
      renameFile tempName fileName)
  putStrLn "These are your TO-DO items:"
  mapM_ putStrLn numberedTasks
remove _ = putStrLn "The remove command takes exactly two arguments"

doesntExist :: String -> [String] -> IO ()
doesntExist command _ = 
  putStrLn $ "The " ++ command ++ " command doesn't exist"
