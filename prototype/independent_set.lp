
task(T) :- task(T,_,_). 
worker(W) :- worker2task(W,_).

% independent sets

edge(T1,T2) :- 
        worker2task(W,T1),
        worker2task(W,T2), 
        T1 < T2,
        task(T1,S1,E1), task(T2,S2,E2),  
        E2 <= S1. 

edge(T1,T2) :- 
        worker2task(W,T1),
        worker2task(W,T2), 
        T1 < T2,
        task(T1,S1,E1), task(T2,S2,E2),  
        E1 <= S2. 

{ set(T) } :- task(T). 

:- edge(T1,T2), set(T1), set(T2). 

#hide. 
#show set/1. 
#show total/1. 
#maximize[ set(_) ].  


% combine with simple lower bound
% TRIVIAL LOWER BOUNDS

intersects(T1,T2) :- 
    task(T1,S1,_),  
    S2 <= S1, 
    S1 < E2, 
    T1 != T2, 
    task(T2,S2,E2). 

max_single(T1,N) :-
    task(T1),
    N = #count{intersects(T1,_)}. 

lower_bound(N) :- 
    N = #max[max_single(_,M)=M]. 

:-  { set(_) } N, lower_bound(N). 

#show lower_bound/1. 
