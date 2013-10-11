 
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
