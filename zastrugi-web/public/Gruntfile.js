'use strict';

module.exports = function(grunt) {
  // load all grunt tasks
  //require('matchdep').filterDev('grunt-*').forEach(grunt.loadNpmTasks);
  require('load-grunt-tasks')(grunt);
  require('time-grunt')(grunt);

  // Project configuration.
  grunt.initConfig({

    meta: {
      version: '1.0.0',
      banner: '/*! TOPOS - v<%= meta.version %> - ' +
        '<%= grunt.template.today("yyyy-mm-dd") %>\n' +
        '* Copyright (c) <%= grunt.template.today("yyyy") %> ' +
        'Ryan Alfonso Roemmich; Licensed MIT */',
      publicDir: '.',
      distDir: 'dist',
      templateDir: '../templates',
      templateDistDir: '../templates/dist',
    },

    clean: {
      dev: {
        src: [
            '<%= meta.publicDir %>/css/',
            '<%= meta.publicDir %>/js/compiled',
        ]
      },
      dist: {
        src: [
            '<%= meta.distDir %>/css/',
            '<%= meta.distDir %>/js/',
        ]
      }
    },

    uglify: {
      options: {
        banner: '/* hey you guise */',
        mangle: {
          toplevel: true,
        },
        compress: {
          global_defs: {
            "DEBUG": false
          },
          dead_code: true
        }
      },
    },

    jshint: {
      options: {
        node: true,
        browser: true,
        esnext: true,
        bitwise: true,
        camelcase: true,
        curly: true,
        eqeqeq: true,
        immed: true,
        indent: 2,
        latedef: true,
        newcap: true,
        noarg: true,
        quotmark: 'single',
        regexp: true,
        undef: true,
        unused: 'vars',
        strict: true,
        trailing: true,
        smarttabs: true,
        jquery: true,
        globals: {
            'topos': true,
            'google': true,
            '_': true,
            'moment': true,
            'Backbone': true,
        }
      },

      dist: {
        src: [
          '<%= meta.publicDir %>/js/*.jsx',
          '<%= meta.publicDir %>/js/**/*.jsx'
        ]
      }
    },

    useminPrepare: {
      html: '<%= meta.templateDistDir %>/index.html',
      options: {
        dest: '<%= meta.publicDir %>',
        root: '<%= meta.publicDir %>',
        flow: {
          steps: {
            js: ['concat', 'uglifyjs'],
          },
          post: {
            js: [{
              name: 'uglifyjs',
              createConfig: function (context, block) {
                context.options = {
                  mangle: true,
                };
              }
            }]
          }
        }
      }
    },

    usemin: {
        html: ['<%= meta.templateDistDir %>/index.html'],
        options: {
            assetsDirs: ['<%= meta.publicDir %>']
        }
    },

    copy: {
        dist: {
            files: [{
                expand: true,
                dot: true,
                cwd: '<%= meta.templateDir %>',
                dest: '<%= meta.templateDistDir %>',
                src: [
                    'index.html'
                ]
            }]
        }
    },

    compass: {
      dev: {
        options: {
          config: '<%= meta.publicDir %>/config.rb',
          basePath: '<%= meta.publicDir %>',
          cssDir: '<%= meta.publicDir %>/css/',
        }
      },

      dist: {
        options: {
          config: '<%= meta.publicDir %>/config.rb',
          basePath: '<%= meta.publicDir %>',
          cssDir: '<%= meta.distDir %>/css/',
        }
      }
    },

    react: {
      files: {
        expand: true,
        cwd: '<%= meta.publicDir %>/js/',
        src: ['*.jsx', '**/*.jsx'],
        dest: '<%= meta.publicDir %>/js/compiled/',
        ext: '.js'
      }
    },

    dowatch: {
      options: {
        livereload: true
      },
      scripts: {
        files: ['<%= meta.publicDir %>/js/*.js',
                '<%= meta.publicDir %>/js/models/*.js',
                '<%= meta.publicDir %>/js/components/*.js',
                '<%= meta.publicDir %>/js/*.jsx',
                '<%= meta.publicDir %>/js/**/*.jsx'],
        //tasks: ['react', 'jshint'],
        tasks: ['react'],
      },
      styles: {
        files: '<%= meta.publicDir %>/sass/{,*/}*.scss',
        tasks: 'compass:dev',
        options: {
          livereload: false
        }
      },
      compiled_css: {
        files: '<%= meta.publicDir %>/css/{,*/}*.css',
        tasks: []
      }
    },

    compress: {
      dist: {
        options: {
          mode: 'gzip'
        },
        files: [
          {
            expand: true,
            src: '<%= meta.publicDir %>/js/dist/*.min.js',
            dest: '.',
            ext: '.js.gz'
          },
          {
            expand: true,
            src: ['<%= meta.publicDir %>/css/dist/*.min.css'],
            dest: '.',
            ext: '.css.gz'
          }
        ]
      }
    },

    filerev: {
      options: {
        encoding: 'utf8',
        algorithm: 'md5',
        length: 8
      },
      js: {
        src: '<%= meta.distDir %>/js/*.js'
      },
      css: {
        src: '<%= meta.distDir %>/css/*.css'
      }
    },
  });

  grunt.renameTask('watch', 'dowatch');
  grunt.registerTask('watch', ['compass:dev', 'dowatch']);

  grunt.registerTask('dist', [
    'clean:dist',
    'compass:dist',
    'copy',
    'react',
    'useminPrepare',
    'concat:generated',
    'uglify:generated',
    'filerev',
    'usemin'
  ]);

  // Default task.
  grunt.registerTask('default', ['react', 'jshint', 'compass:dev']);
  grunt.registerTask('build', ['react', 'jshint', 'compass:dev']);
};

