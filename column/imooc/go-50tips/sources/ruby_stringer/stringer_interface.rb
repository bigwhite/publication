#!/usr/bin/env ruby
# -*- coding: UTF-8 -*-

def concat(a, b)
  unless a.respond_to?(:string) && b.respond_to?(:string)
    raise ArgumentError, "无效参数"
  end
  a.string() + b.string()
end

class Stringer
  def string
    raise NotImplementedError
  end
end

class Foo 
  def initialize(i)
      @i=i
  end
  def string
    @i.to_s
  end
end

class Bar
  def initialize(s)
      @s=s
  end
  def string
    @s
  end
end

f = Foo.new(5)
b = Bar.new("hello")
puts concat(b, f)
