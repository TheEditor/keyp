# bash completion for keyp
# Source this file to enable tab completion for keyp commands

_keyp_completion() {
  local cur prev opts base
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"

  # All available commands
  local commands="init set get list delete rename copy export import sync"

  # Sync subcommands
  local sync_commands="init push pull status config"

  # Flags for various commands
  local global_flags="--help --version"
  local get_flags="--stdout --no-clear --timeout"
  local list_flags="--search --count"
  local delete_flags="-f --force"
  local export_flags="--plain --stdout"
  local import_flags="--replace --dry-run"
  local sync_init_flags="-a -c --auto-push --auto-commit"
  local sync_pull_flags="-s -c --strategy --auto-resolve"
  local sync_config_flags="--auto-push --auto-commit"

  # Handle command completion
  case "${COMP_WORDS[1]}" in
    get|delete|rename|copy)
      if [[ ${COMP_CWORD} -eq 2 ]]; then
        # Complete with secret names
        local secrets=$(keyp list 2>/dev/null | grep -oP '  • \K.*' | sort)
        COMPREPLY=( $(compgen -W "${secrets}" -- ${cur}) )
      fi

      # Add flags for get command
      if [[ "${COMP_WORDS[1]}" == "get" ]] && [[ ${cur} == -* ]]; then
        COMPREPLY=( $(compgen -W "${get_flags}" -- ${cur}) )
      fi

      # Add flags for delete command
      if [[ "${COMP_WORDS[1]}" == "delete" ]] && [[ ${cur} == -* ]]; then
        COMPREPLY=( $(compgen -W "${delete_flags}" -- ${cur}) )
      fi
      ;;

    list)
      if [[ ${cur} == -* ]]; then
        COMPREPLY=( $(compgen -W "${list_flags}" -- ${cur}) )
      fi
      ;;

    export)
      if [[ ${cur} == -* ]]; then
        COMPREPLY=( $(compgen -W "${export_flags}" -- ${cur}) )
      fi
      ;;

    import)
      if [[ ${cur} == -* ]]; then
        COMPREPLY=( $(compgen -W "${import_flags}" -- ${cur}) )
      else
        # Complete with filenames
        COMPREPLY=( $(compgen -f -- ${cur}) )
      fi
      ;;

    sync)
      if [[ ${COMP_CWORD} -eq 2 ]]; then
        COMPREPLY=( $(compgen -W "${sync_commands}" -- ${cur}) )
      elif [[ "${COMP_WORDS[2]}" == "init" ]] && [[ ${cur} == -* ]]; then
        COMPREPLY=( $(compgen -W "${sync_init_flags}" -- ${cur}) )
      elif [[ "${COMP_WORDS[2]}" == "pull" ]] && [[ ${cur} == -* ]]; then
        COMPREPLY=( $(compgen -W "${sync_pull_flags}" -- ${cur}) )
      elif [[ "${COMP_WORDS[2]}" == "config" ]] && [[ ${cur} == -* ]]; then
        COMPREPLY=( $(compgen -W "${sync_config_flags}" -- ${cur}) )
      fi
      ;;

    *)
      if [[ ${cur} == -* ]]; then
        COMPREPLY=( $(compgen -W "${global_flags}" -- ${cur}) )
      else
        COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
      fi
      ;;
  esac

  return 0
}

complete -o bashdefault -o default -o nospace -F _keyp_completion keyp
