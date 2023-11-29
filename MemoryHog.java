/*
 * App that consumes more and more memory until it crashes.
 *
 * @author Torstein Krause Johansen
 */
import java.util.ArrayList;
import java.util.List;
import java.lang.ProcessHandle;

public class MemoryHog {
  public static void main(String[] args) throws Exception {
    long currentProcessId = ProcessHandle.current().pid();
    System.out.println(currentProcessId);

    Thread.sleep(10000L);
    List<Object> list = new ArrayList<Object>();
    while (true) {
      list.add(new Object());
    }
  }
}
