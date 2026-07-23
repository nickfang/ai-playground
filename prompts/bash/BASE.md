When using printf, don't directly print a variable.  use a format like printf '%s\n' "$s"
In a bash script try not to create a subshell if possible.  i.e. use brace expansion instead of seq.  
