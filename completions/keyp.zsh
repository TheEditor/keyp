#compdef keyp
# zsh completion for keyp

_keyp() {
  local -a subcommands
  subcommands=(
    'init:Initialize a new vault'
    'set:Store a secret in the vault'
    'get:Retrieve a secret from the vault'
    'list:List all secrets in the vault'
    'delete:Delete a secret from the vault'
    'rename:Rename a secret'
    'copy:Copy a secret to a new name'
    'export:Export secrets to file'
    'import:Import secrets from file'
    'sync:Synchronize vault with Git remote'
  )

  local -a sync_subcommands
  sync_subcommands=(
    'init:Initialize Git sync with remote'
    'push:Push encrypted vault to remote'
    'pull:Pull vault from remote'
    'status:Show Git sync status'
    'config:Configure Git sync settings'
  )

  local context state line

  _arguments \
    '(-h --help)'{-h,--help}'[Show help]' \
    '(-v --version)'{-v,--version}'[Show version]' \
    '1: :->commands' \
    '*:: :->args'

  case $state in
    commands)
      _describe 'commands' subcommands
      ;;
    args)
      case $words[2] in
        get)
          _arguments \
            '(--stdout)--stdout[Print to stdout instead of clipboard]' \
            '(--no-clear)--no-clear[Do not auto-clear clipboard]' \
            '(--timeout)--timeout[Auto-clear timeout in seconds]:seconds' \
            ':secret-name:_keyp_secrets'
          ;;
        set)
          _arguments \
            ':secret-name:' \
            ':secret-value:'
          ;;
        delete|rm)
          _arguments \
            '(-f --force)'{-f,--force}'[Skip confirmation]' \
            ':secret-name:_keyp_secrets'
          ;;
        rename)
          _arguments \
            ':old-name:_keyp_secrets' \
            ':new-name:'
          ;;
        copy)
          _arguments \
            ':source:_keyp_secrets' \
            ':dest:'
          ;;
        list)
          _arguments \
            '(--search)--search[Search by pattern]:pattern' \
            '(--count)--count[Show only count]'
          ;;
        export)
          _arguments \
            '(--plain)--plain[Export as plaintext JSON]' \
            '(--stdout)--stdout[Print to stdout]' \
            ':output-file:_files'
          ;;
        import)
          _arguments \
            '(--replace)--replace[Replace all existing secrets]' \
            '(--dry-run)--dry-run[Preview without importing]' \
            ':input-file:_files'
          ;;
        sync)
          _arguments '1: :->sync_commands'
          case $state in
            sync_commands)
              _describe 'sync commands' sync_subcommands
              ;;
          esac

          case $words[3] in
            init)
              _arguments \
                '(-a --auto-push)'{-a,--auto-push}'[Enable auto-push]' \
                '(-c --auto-commit)'{-c,--auto-commit}'[Enable auto-commit]' \
                ':remote-url:'
              ;;
            pull)
              _arguments \
                '(-s --strategy)'{-s,--strategy}'[Conflict resolution strategy]:strategy:(keep-local keep-remote)' \
                '(--auto-resolve)--auto-resolve[Auto resolve conflicts]'
              ;;
            config)
              _arguments \
                '(--auto-push)--auto-push[Enable/disable auto-push]:enabled:(true false)' \
                '(--auto-commit)--auto-commit[Enable/disable auto-commit]:enabled:(true false)'
              ;;
          esac
          ;;
      esac
      ;;
  esac
}

# Helper to complete secret names
_keyp_secrets() {
  local secrets
  secrets=(${(f)"$(keyp list 2>/dev/null | grep -oP '  • \K.*')"})
  _values 'secrets' $secrets
}

_keyp "$@"
