---
hide:
  - toc
title: "vkv completion powershell"
---
## vkv completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	vkv completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
vkv completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [vkv completion](vkv_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 30-Apr-2024