#lang racket
(for-each (lambda (arg)
        (printf "~a~n" arg))
        (vector->list (current-command-line-arguments)))
