The Untitled James Condron Encypted Messages Service
==

Given a passworded RSA file, a rabbitmq and some friends, send encrypted messages in a quick and, hopefully secure manner.

This tool utilises RSA keypairs and rabbit to do stuff.

Why RSA keypairs?
--

Easy to implement, nice and secure(as far as anyone can say) and familiar enough. It is comparatively slow, of course, but that shouldn't matter for simple messaging.

Why Rabbit?
--

Purely for the fact messaging is quick, sequential and we know how to scale it.

Licence
--

The MIT License (MIT)
Copyright (c) 2016 jspc

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
