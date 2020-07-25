if exists('g:loaded_project')
    finish
endif
let g:loaded_project = 1
let s:plugin_root = fnamemodify(resolve(expand('<sfile>:p')), ':h:h')

function! s:StartPlugin(host) abort
  return jobstart([s:plugin_root.'/bin/project', 'neovim'], {'rpc': v:true})
endfunction

call remote#host#Register('project', 'x', function('s:StartPlugin'))
call remote#host#RegisterPlugin('project', '0', [
    \   {'type': 'function', 'name': 'Alternate', 'sync': 1, 'opts': {}},
    \   {'type': 'function', 'name': 'CompleteRelatedKey', 'sync': 1, 'opts': {}},
    \ ])


command! A call Alternate('tabe')


