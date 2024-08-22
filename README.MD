# Vector

An abstract vector and based on it data structures.

## Motivation

The golang slices are great, and in most cases they are complete enoght to implement
any logic related to collections. But, in some cases, a phylosopy behind a slices
not works or produces a lot of copypaste. This package is an attempt to
implement some "C++ STL like" data structures and reuse it between logic where
they are applicable.

## Conclustrions

Do not use this package if you are not sure. Just use a slices. If you decided to use
or extend it, please read 'Conventions'.

## Conventions

* All data types are based on VectorImpl generic.
* All API are separated to 2 peases: thread safe and not.
* Every thread safe method must have an not thread safe analog.
* In most cases, the thread safe method must just perforom lock and call not thread safe analog.
* A not thread safe API is the package private API.
* The package pablic APIs must be a thread safe.
* Sure there are an exceptions from that rules, but they have to be clearly motivated somehow.

# Maintainers:  

- Dmytro Iakovliev <d.iakovliev@dtpoland.com>
- CS23 team <cs23@dtpoland.com>