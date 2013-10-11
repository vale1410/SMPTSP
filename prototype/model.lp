#hide. 
#show assign/2. 

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

%is_start(W,ST) :- assign(T,W), task(T,ST,_). 
%is_end(W,ET) :- assign(T,W), task(T,_,ET). 
%
%order_start(W,C1,C2) :- task(T1,C1,_), task(T2,C2,_), worker2task(W,T1), worker2task(W,T2),
%                   C1 < C2 , C <= C1 : task(T,C,_) : worker2task(W,T) : C < C2.
%order_end(W,C1,C2) :- task(T1,_,C1), task(T2,_,C2), worker2task(W,T1), worker2task(W,T2),
%                   C1 < C2 , C <= C1 : task(T,_,C) : worker2task(W,T) : C < C2.

%is_safe_start(W,ST). 
%is_safe_end(W,ST). 
%
%order_st(ST1,ST2). 
%
%next_et(ST,ET). 
%
%is_safe_start(W,S2) :-  
%

%https://sites.google.com/site/ptsplib/smptsp/home
