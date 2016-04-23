@test "dexec is present" {
  run dexec -v
  [ "$status" -eq 0 ]
  grep -Ee '^dexec [[:digit:]]+.[[:digit:]]+.[[:digit:]]+(-SNAPSHOT)?$' <<<"$output"
}

function run_standard_tests() {
  pushd $BATS_TEST_DIRNAME/fixtures/$1 >/dev/null
  run dexec [Hh]ello[Ww]orld*
  [ "$status" -eq 0 ]
  [ "$output" = "hello world" ]

  run dexec [Uu]nicode*
  [ "$status" -eq 0 ]
  [ "$output" = "hello unicode 👾" ]

  run dexec [Ee]cho[Cc]hamber* -a hello -a world -a 'test with spaces'
  [ "$status" -eq 0 ]
  [ "${lines[0]}" = "hello" ]
  [ "${lines[1]}" = "world" ]
  [ "${lines[2]}" = "test with spaces" ]

  run ./[Ss]hebang*
  [ "$status" -eq 0 ]
  [ "$output" = "hello world" ]
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
  skip
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "lisp" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "node" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "objc" {
  skip
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
  skip
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "rust" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}

@test "scala" {
  run_standard_tests $BATS_TEST_DESCRIPTION
}
