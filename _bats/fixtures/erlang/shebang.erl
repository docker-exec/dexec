#!/usr/bin/env dexec
%% -*- erlang -*-
main([]) ->
    io:setopts([{encoding, unicode}]),
    io:fwrite("hello world\n").
