struct Ticket {
    name: String,
}

impl Drop for Ticket {
    fn drop(&mut self) {
        println!("[Drop] {} 的门票被销毁了", self.name);
    }
}

fn take_ownership(ticket: Ticket) {
    println!("[take_ownership] 拿到了 {} 的门票，进场", ticket.name);
}

fn give_back(ticket: Ticket) -> Ticket {
    println!("[give_back] 检查一下 {} 的门票，还给你", ticket.name);
    ticket
}

fn main() {
    let t1 = Ticket {
        name: String::from("张三"),
    };
    take_ownership(t1);
    println!("take_ownership 调用结束");

    let t2 = Ticket {
        name: String::from("李四"),
    };
    let t2 = give_back(t2);
    println!("t2 还在我手里：{}", t2.name);

    {
        let t3 = Ticket {
            name: String::from("王五"),
        };
        println!("t3 在内层作用域：{}", t3.name);
    }
    println!("离开内层作用域之后");
}
