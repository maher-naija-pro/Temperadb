module.exports = {
    extends: ['@commitlint/config-conventional'],
    rules: {
        // Type must be one of the conventional types
        'type-enum': [
            2,
            'always',
            [
                'feat',     // New feature
                'fix',      // Bug fix
                'docs',     // Documentation changes
                'style',    // Code style changes (formatting, missing semicolons, etc)
                'refactor', // Code refactoring
                'perf',     // Performance improvements
                'test',     // Adding or updating tests
                'chore',    // Maintenance tasks
                'ci',       // CI/CD changes
                'build',    // Build system changes
                'revert',   // Revert previous commit
                'security', // Security fixes
                'deps',     // Dependency updates
                'wip',      // Work in progress
                'hotfix',   // Critical hotfix
                'release'   // Release commits
            ]
        ],

        // Subject must not be empty
        'subject-empty': [2, 'never'],

        // Subject must not end with a period
        'subject-full-stop': [2, 'never', '.'],

        // Subject must be in sentence case
        'subject-case': [2, 'always', 'sentence-case'],

        // Subject must not be longer than 72 characters
        'subject-max-length': [2, 'always', 72],

        // Body must not be longer than 100 characters per line
        'body-max-line-length': [2, 'always', 100],

        // Footer must not be longer than 100 characters per line
        'footer-max-line-length': [2, 'always', 100],

        // Type must be lowercase
        'type-case': [2, 'always', 'lower'],

        // Type must not be empty
        'type-empty': [2, 'never'],

        // Scope must be lowercase
        'scope-case': [2, 'always', 'lower'],

        // Scope must be one of the allowed scopes
        'scope-enum': [
            2,
            'always',
            [
                'api',        // API changes
                'auth',       // Authentication changes
                'build',      // Build system
                'ci',         // CI/CD
                'cli',        // Command line interface
                'config',     // Configuration changes
                'core',       // Core functionality
                'db',         // Database changes
                'docs',       // Documentation
                'feat',       // Features
                'fix',        // Bug fixes
                'lint',       // Linting rules
                'perf',       // Performance
                'refactor',   // Refactoring
                'security',   // Security
                'storage',    // Storage layer
                'test',       // Testing
                'ui',         // User interface
                'utils',      // Utilities
                'web',        // Web interface
                'deps',       // Dependencies
                'release'     // Release management
            ]
        ],

        // Breaking changes must be indicated
        'breaking-enum': [2, 'always', ['!', 'BREAKING CHANGE']],

        // References to issues must be valid
        'references-empty': [1, 'never'],

        // Signed-off-by must be present for certain types
        'signed-off-by': [1, 'always', 'Signed-off-by:'],

        // Custom rules for TimeSeriesDB project
        'header-max-length': [2, 'always', 100],
        'body-leading-blank': [2, 'always'],
        'footer-leading-blank': [2, 'always'],

        // Enforce conventional commit format
        'header-format': [2, 'always', 'lower'],

        // Require scope for certain types
        'scope-required': [
            2,
            'always',
            ['feat', 'fix', 'refactor', 'perf', 'test', 'docs']
        ],

        // Enforce ticket references for certain types
        'ticket-required': [
            1,
            'always',
            ['feat', 'fix', 'refactor', 'perf']
        ]
    },

    // Custom parser for ticket references
    parserPreset: {
        parserOpts: {
            issuePrefixes: ['TSDB-', 'FIXES-', 'CLOSES-', 'RELATES-']
        }
    },

    // Helpful messages for common violations
    helpUrl: 'https://github.com/conventional-changelog/commitlint/#what-is-commitlint',

    // Ignore certain commits (like merge commits)
    ignores: [
        '^Merge branch',
        '^Merge pull request',
        '^Revert',
        '^WIP:',
        '^wip:',
        '^fixup!',
        '^squash!'
    ]
};
