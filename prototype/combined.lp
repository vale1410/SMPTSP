task(T) :- task(T,_,_). 
worker(W) :- worker2task(W,_).

1 { assign(T,W) : worker2task(W,T) } 1 :- task(T). 


%%% SIMPLE APPROACH %%%%

:-  assign(T1,W), assign(T2,W), T1 != T2, 
    task(T1,S1,E1), task(T2,S2,E2),  
    S1 <= S2, S2 < E1.  


% LOWER BOUND!

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

%% TRIVIAL LOWER BOUNDS
%
%intersects(T1,T2) :- 
%    task(T1,S1,_),  
%    S2 <= S1, 
%    S1 < E2, 
%    T1 != T2, 
%    task(T2,S2,E2). 
%
%max_single(T1,N) :-
%    task(T1),
%    N = #count{intersects(T1,_)}. 
%
%lower_bound(N) :- 
%    N = #max[max_single(_,M)=M]. 
%
%:-  { set(_) } N, lower_bound(N). 
%
%
bound(N) :- N = #count{used(_)}. 
#show bound/1. 

% LOWER AND UPPER BOUND WORK TOGETHER...

#hide. 

used(W) :- assign(_,W). 

% CHANNELING

% direction 1: if worker is selected then one of the sets have to be selected!
:- used(W), not set(T) : worker2task(W,T). 

% direction 2: 
:- set(T), not used(W) : worker2task(W,T). 

:- not 0 [used(_) = 1, set(_) = -1] 0. 
%:- 7 [used(_) = 1, set(_) = -1]. 

%#minimize[used(_) = 1, set(_) = -1]. 

%#show count/1. 
%count(T) :- T = [used(_)].

