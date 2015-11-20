[<EntryPoint>]
let main args =
    args |> Array.iter (fun x -> printfn "%s" x)
    0
