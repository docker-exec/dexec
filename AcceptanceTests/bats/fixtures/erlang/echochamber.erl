%% -*- erlang -*-

format_list(L) ->
    fnl(L).

fnl([H]) ->
    io:format("~s~n", [H]);

fnl([H|T]) ->
    io:format("~s~n", [H]),
    fnl(T);

fnl([]) ->
    ok.

main(Args) ->
    io:setopts([{encoding, unicode}]),
    format_list(Args).
