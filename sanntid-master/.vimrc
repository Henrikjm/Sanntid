syntax on
"color xoria256
set autoindent
set ts=4 st=4 sw=4 noexpandtab
" Omnicomplete
set ofu=syntaxcomplete#Complete
" Set visible tab and <CR> symbols
set listchars=tab:▸\ ,eol:¬
set nocompatible
set foldmethod=syntax
set foldnestmax=1

" Mappings
nmap <F2> :tabprev <CR>
nmap <F3> :tabnext <CR>
nmap <F4> :tabnew <CR>
nmap <F5> :set hlsearch! <CR>
nmap <F6> :set number! <CR>
nmap <F7> :set list! <CR>
