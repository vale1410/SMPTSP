task(T) :- task(T,_,_). 
worker(W) :- worker2task(W,_).

% find simple lower bound

intersects(T1,T2) :- 
    task(T1,S1,_),  
    S2 <= S1, 
    S1 < E2, 
    T1 != T2, 
    task(T2,S2,E2). 

max_single(T1,N) :-
    task(T1),
    N = #count{intersects(T1,_)}. 

lower_bound(N+1) :- 
    N = #max[max_single(_,M)=M]. 


#hide. 
#show lower_bound/1. 
