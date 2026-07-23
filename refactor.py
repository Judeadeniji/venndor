import re
import os

with open('internal/cli/cli.go', 'r') as f:
    content = f.read()

# We want to match:
# var nameCmd = &cobra.Command{
#   Use: "...",
#   Short: "...",
#   Args: ...,
#   Run: func(cmd *cobra.Command, args []string) { ... },
# }

def replacer(match):
    cmd_name = match.group(1) # e.g. addCmd
    func_name = 'run' + cmd_name[0].upper() + cmd_name[1:-3] # e.g. runAdd
    
    body = match.group(2)
    
    # We replace the Run inline function with the named function.
    new_cmd_decl = match.group(0).replace(
        'Run: func(cmd *cobra.Command, args []string) {\n' + body + '\n\t},',
        f'Run: {func_name},'
    )
    
    # Prepend the new function definition
    new_func = f'func {func_name}(cmd *cobra.Command, args []string) {{\n{body}\n}}\n\n'
    
    return new_func + new_cmd_decl

# Regex to match var XxxCmd = &cobra.Command{ ... Run: func(...) { BODY }, ... }
# We need to capture the cmd name, and the body of the Run func.
pattern = re.compile(
    r'var (\w+Cmd) = &cobra\.Command\{[^{]+'
    r'Run:\s*func\(cmd \*cobra\.Command, args \[\]string\) \{\n(.*?)\n\t\},'
    r'\n\}',
    re.DOTALL
)

new_content = pattern.sub(replacer, content)

with open('internal/cli/cli.go', 'w') as f:
    f.write(new_content)

print("Refactored successfully")
