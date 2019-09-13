@test "dexec is present" {
  run dexec -v
  [ "$status" -eq 0 ]
  grep -Ee '^dexec [[:digit:]]+.[[:digit:]]+.[[:digit:]]+(-SNAPSHOT)?$' <<<"$output"
}

function run_standard_tests() {
  pushd $BATS_TEST_DIRNAME/fixtures/$1 >/dev/null
  run dexec [Hh]ello[Ww]orld*
  [ "$status" -eq 0 ]
  [ "$output" = $'hello world\r' ]

  run dexec [Uu]nicode*
  [ "$status" -eq 0 ]
  [ "$output" = $'hello unicode 👾\r' ]

  run dexec [Ee]cho[Cc]hamber* -a hello -a world -a 'test with spaces'
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = $'hello\r' ]
  [ "${lines[1]}" = $'world\r' ]
  [ "${lines[2]}" = $'test with spaces\r' ]

  run ./[Ss]hebang*
  [ "$status" -eq 0 ]
  [ "$output" = $'hello world\r' ]
  popd >/dev/null
}

@test "bash" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "c" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "clojure" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "coffee" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "cpp" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "csharp" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "d" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "erlang" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "fsharp" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "go" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "groovy" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "haskell" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "java" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "lisp" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "node" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "objc" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "ocaml" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "perl" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "php" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "python" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "racket" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "ruby" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "rust" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "scala" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}
