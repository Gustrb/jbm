interface Animal {
    void eat();
}

interface Mammal extends Animal {
    void move();
}

public class Person implements Mammal {

    public void eat() {
        System.out.println("Eating");
    }

    public void move() {
        System.out.println("Moving");
    }
}
