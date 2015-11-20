open Printf

let () =
    for i = 1 to Array.length Sys.argv - 1 do
      printf "%s\n" Sys.argv.(i)
    done;;
