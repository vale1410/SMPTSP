#hide. 

task(T) :- task(T,_,_). 
worker(W) :- worker2task(W,_).

1 { assign(T,W) : worker2task(W,T) } 1 :- task(T). 

used(W) :- assign(_,W). 

#minimise{used(_)}.

%%% SIMPLE APPROACH %%%%

:-  assign(T1,W), assign(T2,W), T1 != T2, 
    task(T1,S1,E1), task(T2,S2,E2),  
    S1 <= S2, S2 < E1.  

%%% MORE COMPLEX %%%%
%  when is a start time connected backward to another start time?
%  post S1 <- S2, when you should not start at S2, then also not at S1 ?
%  if all tasks that forbid S2, also forbid S1!
%  todo: collect tasks that forbid S. 
%  if all tasks that forbid S2, and next(S1,S2), then introduce the
%  implication.

% T1 blocks start time S2
block(W,T1,S2) :-  
    T1 != T2, 
    worker2task(W,T1), 
    worker2task(W,T2),
    task(T1,S1,E1), task(T2,S2,E2),  
    S1 <= S2, S2 < E1.  

succ(W,S1,S2) :- 
    task(T1,S1,_), task(T2,S2,_), 
    worker2task(W,T1), worker2task(W,T2),
    S1 < S2, 
    S <= S1 : task(T,S,_) : worker2task(W,T) : S < S2. 

block_succ(W,S1,S2) :- 
    succ(W,S1,S2), 
    block(W,T,S1) : block(W,T,S2). 

block_max(W,S) :- 
    block_succ(W,_,S), 
    not block_succ(


start(W,S) :- assign(T,W), task(T,S,_). 

dont_start(W,S) :-  assign(T,W), blocks_max(W,T,S). %of everything only take the largest

dont_start(W,S1) :-  dont_start(W,S2), blocks_succ(S2,S1). 
