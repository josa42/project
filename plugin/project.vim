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
    \   {'type': 'function', 'name': 'ProjectOpen', 'sync': 1, 'opts': {}},
    \   {'type': 'function', 'name': 'CompleteRelatedKey', 'sync': 1, 'opts': {}},
    \   {'type': 'function', 'name': 'CompleteKey', 'sync': 1, 'opts': {}},
    \   {'type': 'function', 'name': 'CompleteOpen', 'sync': 1, 'opts': {}},
    \ ])

" function! CompleteAlternate(a,b,c)
"   return ['window', 'tab', 'tab!', 'split', 'vsplit']
" endfunction
" command! -nargs=? -complete=customlist,CompleteAlternate A call Alternate(<args>)

if !exists('g:project_default_command')
    let g:project_default_command = 'tab'
endif


command! -nargs=0 -bang A call Alternate(g:project_default_command . '<bang>')

command! -nargs=0 -bang AW call Alternate('window<bang>')
command! -nargs=0 -bang AT call Alternate('tab<bang>')
command! -nargs=0 -bang AS call Alternate('split<bang>')
command! -nargs=0 -bang AV call Alternate('vsplit<bang>')

nmap <silent> <space>o :call ProjectOpen(g:project_default_command)<cr>
nmap <silent> <space>a :call Alternate(g:project_default_command)<cr>
