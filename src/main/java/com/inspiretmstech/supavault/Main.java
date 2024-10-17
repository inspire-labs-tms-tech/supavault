package com.inspiretmstech.supavault;

import ch.qos.logback.classic.Level;
import ch.qos.logback.classic.Logger;
import com.inspiretmstech.supavault.bases.Executor;
import com.inspiretmstech.supavault.constants.Version;
import com.inspiretmstech.supavault.env.RuntimeEnvironment;
import org.slf4j.LoggerFactory;
import picocli.CommandLine;

import java.util.Objects;

@CommandLine.Command(
        name = "supavault",
        subcommands = {},
        versionProvider = Version.class,
        mixinStandardHelpOptions = true
)
class Main implements Runnable {

    private static final Logger logger = (Logger) LoggerFactory.getLogger(Logger.ROOT_LOGGER_NAME);
    private static final org.slf4j.Logger log = LoggerFactory.getLogger(Main.class);

    static {
        System.setProperty("slf4j.internal.verbosity", "OFF");
        ((Logger) LoggerFactory.getLogger(Logger.ROOT_LOGGER_NAME)).setLevel(Level.ERROR);
    }

    @CommandLine.Option(names = {"-v", "--version"}, versionHelp = true,
            description = "print version information and exit")
    boolean versionRequested;

    private static int execute(Executor<Integer> executor) {
        try {
            return executor.execute();
        } catch (RuntimeException e) {
            logger.error(e.getMessage());
        } catch (Exception e) {
            logger.error("A critical error occurred: {}", e.getMessage());
        }
        return 1;
    }

    public static void main(String[] args) {
        System.exit(Main.execute(() -> new CommandLine(new Main()).execute(args)));
    }

    @CommandLine.Option(names = {"-H", "--db-host"}, description = "host of the Postgres database (e.g., 127.0.0.1 or my.db.example.com)", scope = CommandLine.ScopeType.INHERIT)
    private void dbHost(String dbHost) {
        if (Objects.nonNull(dbHost) && !dbHost.isEmpty())
            RuntimeEnvironment.setDbHost(dbHost);
    }

    @CommandLine.Option(names = {"-P", "--db-port"}, description = "port of the Postgres database (e.g., 54321)", scope = CommandLine.ScopeType.INHERIT)
    private void dbPort(String dbPort) {
        if (Objects.nonNull(dbPort) && !dbPort.isEmpty())
            RuntimeEnvironment.setDbPort(dbPort);
    }

    @CommandLine.Option(names = {"-N", "--db-name"}, description = "name of the Postgres database (e.g., postgres)", scope = CommandLine.ScopeType.INHERIT)
    private void dbName(String dbName) {
        if (Objects.nonNull(dbName) && !dbName.isEmpty())
            RuntimeEnvironment.setDbName(dbName);
    }

    @CommandLine.Option(names = {"-u", "--db-user"}, description = "username for the user to authenticate into the Postgres database with (e.g., postgres)", scope = CommandLine.ScopeType.INHERIT)
    private void dbUser(String dbUser) {
        if (Objects.nonNull(dbUser) && !dbUser.isEmpty())
            RuntimeEnvironment.setDbUser(dbUser);
    }

    @CommandLine.Option(names = {"-p", "--db-pass"}, description = "password for the user to authenticate into the Postgres database with (e.g., my$uper$ecretPa$$word!)", scope = CommandLine.ScopeType.INHERIT)
    private void dbPass(String dbPass) {
        if (Objects.nonNull(dbPass) && !dbPass.isEmpty())
            RuntimeEnvironment.setDbPass(dbPass);
    }

    @CommandLine.Option(names = {"-V", "--verbose"}, description = "show verbose output", scope = CommandLine.ScopeType.INHERIT)
    public void configureLogger(boolean verbose) {
        Logger rootLogger = (Logger) LoggerFactory.getLogger(Logger.ROOT_LOGGER_NAME);
        rootLogger.setLevel(verbose ? Level.INFO : Level.ERROR);
    }

    @Override
    public void run() {
        this.configureLogger(true);
        logger.info("Supavault {}", Version.version);
        logger.info("> Use --help for more information");
        logger.info("> Report issues at: https://github.com/inspire-labs-tms-tech/supavault/issues");
    }
}
