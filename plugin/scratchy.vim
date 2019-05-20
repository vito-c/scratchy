if exists('g:loaded_scratch')
  finish
endif
let g:loaded_scratch = 1

function! s:Requirescratchy(host) abort
  return jobstart(['/Users/vito.cutten/code/personal/scratchy/bin/scratchy'], {'rpc': v:true})
endfunction

call remote#host#Register('scratchy', 'x', function('s:Requirescratchy'))
call remote#host#RegisterPlugin('scratchy', '0', [
\ {'type': 'function', 'name': 'ScratchyRun', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'ScratchyFormat', 'sync': 1, 'opts': {}},
\ {'type': 'function', 'name': 'ScratchySetup', 'sync': 1, 'opts': {}},
\ ])




" " vim:ts=4:sw=4:et
