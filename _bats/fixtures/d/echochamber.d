import std.stdio;

void main(string[] args)
{
  foreach (string arg; args[1..$])
  {
    writeln(arg);
  }
}
