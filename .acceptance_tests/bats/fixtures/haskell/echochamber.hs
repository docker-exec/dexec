import System.Environment

main = do
   args <- getArgs
   mapM_ putStrLn args
